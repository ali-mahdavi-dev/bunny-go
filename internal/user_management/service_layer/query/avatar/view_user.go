package user

// func ViewUser(ctx context.Context, uow internal.UnitOfWorkImp, cache cache.Store, param *ginx.PaginationResult) (*[]entities.User, error) {
// 	user := new([]entities.User)
// 	key := cache.CreateKey("users")

// 	err := cache.Cache(ctx, key, user, time.Second*10, func(ctx context.Context) (any, error) {
// 		return uow.Do(ctx, func(ctx context.Context, tx *gorm.DB) (any, error) {
// 			if uow.User().Model(ctx).Limit(int(param.Limit)).Offset(int(param.Skip)).Order(param.OrderBy.ToSQL()).Find(user).Count(&param.Total).Error != nil {
// 				return nil, cerrors.BadRequest("Operation.CanNot")
// 			}
// 			return user, nil
// 		})
// 	})

// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil

// }
