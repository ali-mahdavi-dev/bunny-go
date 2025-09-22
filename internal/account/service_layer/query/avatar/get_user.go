package user

// func GetUser(ctx context.Context, uow internal.UnitOfWorkImp, id uint, cache cache.Store) (*entities.User, error) {
// 	user := new(entities.User)
// 	key := cache.CreateKey("user", id)
// 	err := cache.Cache(ctx, key, user, time.Second*5, func(ctx context.Context) (any, error) {
// 		return uow.Do(ctx, func(ctx context.Context, tx *gorm.DB) (any, error) {
// 			if uow.User().Model(ctx).Preload("Trades").First(user, id).Error != nil {
// 				return nil, cerrors.BadRequest("Operation.CanNot")
// 			}
// 			return user, nil
// 		})

// 	})
// 	return user, err

// }
