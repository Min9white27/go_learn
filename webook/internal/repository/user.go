package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// SELECT * FROM `users` WHERE `email`=?
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) FindByIdV1(ctx context.Context, uid interface{}) (domain.User, error) {
	//先从 cache 缓存里面找
	//再从 dao 里面找
	//找到了回写 cache 缓存
	u, err := r.dao.FindByIdV1(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(u), nil
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err != nil {
		//	必然有数据
		return u, nil
	}
	// 没有这个数据
	//if err == cache.ErrKeyNotExist {
	//	//	去数据库里面加载
	//}

	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}

	// 使用 Redis 一定会出现数据不一致性，无可避免
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			//	打日志，做监控就行，不需要下面这句。原因是缓存设置失败，问题不大。做监控是为了防止 Redis 崩了
			//	return domain.User{},err
		}
	}()

	return u, err

	// err = io.EOF
	// 要不要去数据库加载？不加载，然而这个 err 只是偶然发生的错误，不太友好；加载，可能会搞崩数据库。

	// 选加载 —— 要做好兜底，万一 Redis 真的崩了。要保护好数据库（面试必选这个）
	// 如何保护数据库，可以做数据库限流。用 ORM 的 middleware ,不能用 Redis

	// 选不加载 —— 用户体验差一些。实际用这个

	// 缓存里面有数据
	// 缓存里面没有数据
	// 缓存出错了，你也不知道有没有数据
}

func (r *UserRepository) UpdateNonZeroFields(ctx context.Context, user domain.User) error {
	return r.dao.UpdateByUid(ctx, r.toEntity(user))
}

func (r *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:              u.Id,
		Email:           u.Email,
		Password:        u.Password,
		Nickname:        u.Nickname,
		Birthday:        time.UnixMilli(u.Birthday),
		PersonalProfile: u.PersonalProfile,
	}
}

func (r *UserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id:              u.Id,
		Email:           u.Email,
		Password:        u.Password,
		Nickname:        u.Nickname,
		Birthday:        u.Birthday.UnixMilli(),
		PersonalProfile: u.PersonalProfile,
	}
}
