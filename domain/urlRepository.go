package domain

type UrlRepository struct {
	//cli *redis.Client
	urls map[string]string
}

func NewUrlRepository() (*UrlRepository, error) {
	//redisHost := os.Getenv("DB_HOST")
	//redisPort := os.Getenv("DB_PORT")
	//
	//if redisHost == "" || redisPort == "" {
	//	return nil, errors.New("database variables not correctly initiated")
	//}
	//
	//cli := redis.NewClient(&redis.Options{
	//	Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	//})
	//
	//return &UrlRepository{cli: cli}, nil

	urls := make(map[string]string)

	return &UrlRepository{urls: urls}, nil
}

func (ur *UrlRepository) Ping() {
	//val, _ := ur.cli.Ping().Result()
	//log.Printf("Redis URL db ping info: %x", val)
}
