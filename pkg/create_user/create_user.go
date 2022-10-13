package create_user

import (
	"errors"
	"fmt"
	"iredmail-create-email-account/pkg/public_error"
	"iredmail-create-email-account/pkg/remote_ssh"
	"net/http"
	"strings"
	"time"
)

var default_timeout = 60 * time.Second

type CreateUserDTO struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type CreateUserService struct {
	ssh_config *remote_ssh.MakeConfig
}

func New(ssh_private_key_path, server, user, port string) CreateUserService {
	return CreateUserService{
		ssh_config: &remote_ssh.MakeConfig{
			User:    user,
			Server:  server,
			KeyPath: ssh_private_key_path,
			Timeout: 60 * time.Second,
			Port:    port,
		},
	}
}

func (srv CreateUserService) Create(dto *CreateUserDTO) error {
	ssh := srv.ssh_config

	// 1. Get current working directory in order to have easy access to iredmail tool scripts
	pwd, err := srv.pwd(ssh)
	if err != nil {
		return err
	}

	// 2. Generate create user SQL script
	err = srv.create_user_sql_script(pwd, dto, ssh)
	if err != nil {
		return err
	}

	// 3. Execute SQL script
	err = srv.execute_create_user_sql_script(ssh)
	if err != nil {
		return err
	}

	// 4. Remove temp file
	err = srv.remove_temp_sql_script(ssh)
	if err != nil {
		return err
	}

	return nil
}

// Get current working directory on server so we can access to iredmail scripts easier.
func (CreateUserService) pwd(ssh *remote_ssh.MakeConfig) (string, error) {
	pwd_output, _, _, err := ssh.Run("pwd", default_timeout)
	if err != nil {
		err = fmt.Errorf("can't run pwd command: %w", err)
		return "", err
	}

	pwd_output = strings.Replace(pwd_output, "\n", "", -1)
	return pwd_output, nil
}

// Generates SQL script to create a new user in database and stores it into temp directory.
func (CreateUserService) create_user_sql_script(
	pwd string,
	dto *CreateUserDTO,
	ssh *remote_ssh.MakeConfig,
) error {
	create_user_bassh_script_command := fmt.Sprintf(
		"bash %s/iRedMail-1.6.1/tools/create_mail_user_SQL.sh %s '%s' > /tmp/user.sql",
		pwd, dto.UserName, dto.Password,
	)
	_, _, _, err := ssh.Run(create_user_bassh_script_command, default_timeout)
	if err != nil {
		err = fmt.Errorf("can't run create_user_sh_script: %w", err)
		return err
	}
	return nil
}

// Execute temp SQL script previous generated.
func (CreateUserService) execute_create_user_sql_script(ssh *remote_ssh.MakeConfig) error {
	execute_sql_command := "sudo -u postgres psql -d vmail < /tmp/user.sql"
	fmt.Println(execute_sql_command)
	_, stderr, _, err := ssh.Run(execute_sql_command, default_timeout)
	if err != nil {
		err = fmt.Errorf("can't run sql_command: %w", err)
		return err
	}

	if stderr != "" {
		err = fmt.Errorf("can't run sql_command: %w", err)
		if strings.Contains(stderr, "duplicate key") {
			publicErr := public_error.New(
				errors.New("user already exists, try another please"),
				err,
				http.StatusBadRequest,
			)
			return publicErr
		}

		return err
	}

	return nil
}

// Clear temp file.
func (CreateUserService) remove_temp_sql_script(ssh *remote_ssh.MakeConfig) error {
	remove_command := "sudo rm -rf /tmp/user.sql"
	fmt.Println(remove_command)
	_, _, _, err := ssh.Run(remove_command, default_timeout)
	if err != nil {
		err = fmt.Errorf("can't run remove_command: %w", err)
		return err
	}
	return nil
}
