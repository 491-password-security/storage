package Repositories

type RepositoryHandler struct {
	RefreshTokenRepository          RefreshTokenRepository
	SecureTokenRepository           SecureTokenRepository
	JwtRepository                   JwtRepository
	UserRepository                  UserRepository
	RequestRecordRepository         RequestRecordRepository
	BackendGenericRepositories      map[string]*BackendGenericRepository
	DeviceModelStatisticsRepository DeviceModelStatisticsRepository
	SensorDataRepository            SensorDataRepository
	SensorDataLocationRepository    SensorDataLocationRepository
	BackendAnalysisRepository		BackendAnalysisRepository
}
