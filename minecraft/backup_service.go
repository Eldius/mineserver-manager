package minecraft

type BackupService interface {
	Backup(path string) (string, error)
}

type backupService struct {
}

func NewBackupService() BackupService {
	return &backupService{}
}

func (s *backupService) Backup(path string) (string, error) {

	return "", nil
}
