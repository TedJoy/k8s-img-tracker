package core

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"git2.gnt-global.com/jlab/gdeploy/img-tracker/config"
	"git2.gnt-global.com/jlab/gdeploy/img-tracker/pkg/logger"
)

// func GetImgDigestPrivate(user, password, registry, image string) string {
// 	loginPrivateRegistry(user, password, registry)
// 	return GetImgDigestPublic(image)
// }

func loginPrivateRegistry(user, password, registry string) {
}

func (e SkopeoUnauthorizedError) Error() string {
	return "Unauthorized, need login"
}

func tryGetHash(image string) (string, error) {
	cmd := exec.Command("skopeo", "inspect", "docker://"+image, "-f={{.Digest}}")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		if strings.Contains(errb.String(), "unauthorized: authentication required") {
			logger.Logger.Debug(err)
			return "", SkopeoUnauthorizedError{}
		}
		logger.Logger.Infof("error: %v", err)
		return "", err
	}

	return strings.TrimSpace(outb.String()), nil
}

func getECRLoginPassword(key, secret, region string) (string, error) {
	cmd := exec.Command("aws", "ecr", "get-login-password")
	cmd.Env = append(os.Environ(), "AWS_ACCESS_KEY_ID="+key, "AWS_SECRET_ACCESS_KEY="+secret, "AWS_DEFAULT_REGION="+region)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return outb.String(), nil
}

func loginECR(registry, key, secret, region string) error {
	password, err := getECRLoginPassword(key, secret, region)
	if err != nil {
		return err
	}
	cmd := exec.Command("skopeo", "login", registry, "--username=AWS", "--password="+password)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func GetImgDigest(image string, hm map[string]string) string {
	if val, ok := hm[image]; ok {
		return val
	}

	hash, err := tryGetHash(image)

	if err != nil {
		logger.Logger.Debugf("%T", err)
		switch err.(type) {
		case SkopeoUnauthorizedError:
			for k, v := range config.MyFileConfig.Registries {
				if strings.HasPrefix(image, k) {
					logger.Logger.Debug(v.Type)
					switch v.Type {
					case "ECR":
						err = loginECR(k, v.Metadata["aws_access_key_id"], v.Metadata["aws_secret_access_key"], v.Metadata["aws_default_region"])

						if err != nil {
							logger.Logger.Infow("Error when authenticating to ECR", "err object", err)
						}

						hash, err = tryGetHash(image)

						if err != nil {
							logger.Logger.Infow("Error even after authenticated", "err object", err)
						}
					default:
						logger.Logger.Infow("Unsupported registry type", "type", v.Type)
					}
					break
				}
			}
		default:
			logger.Logger.Infow("Error", "err object", err)
		}
	}

	hm[image] = hash

	return hash
}
