package in

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/concourse/s3-resource"
	"github.com/concourse/s3-resource/versions"
)

type InCommand struct {
	s3client s3resource.S3Client
}

func NewInCommand(s3client s3resource.S3Client) *InCommand {
	return &InCommand{
		s3client: s3client,
	}
}

func (command *InCommand) Run(destinationDir string, request InRequest) (InResponse, error) {
	err := command.createDirectory(destinationDir)
	if err != nil {
		return InResponse{}, err
	}

	remotePath, err := command.pathToDownload(request)
	if err != nil {
		return InResponse{}, err
	}

	remoteFilename := path.Base(remotePath)
	err = command.downloadFile(
		request.Source.Bucket,
		remotePath,
		destinationDir,
	)
	if err != nil {
		return InResponse{}, err
	}

	err = command.writeURLFile(
		request.Source.Bucket,
		remotePath,
		request.Source.Private,
		destinationDir,
	)
	if err != nil {
		return InResponse{}, err
	}

	return InResponse{
		Version: s3resource.Version{
			Path: remotePath,
		},
		Metadata: []s3resource.MetadataPair{
			s3resource.MetadataPair{
				Name:  "filename",
				Value: remoteFilename,
			},
		},
	}, nil
}

func (command *InCommand) pathToDownload(request InRequest) (string, error) {
	if request.Version.Path == "" {
		extractions := versions.GetBucketFileVersions(command.s3client, request.Source)

		if len(extractions) == 0 {
			return "", errors.New("no extractions could be found - is your regexp correct?")
		}

		lastExtraction := extractions[len(extractions)-1]
		return lastExtraction.Path, nil
	}

	return request.Version.Path, nil
}

func (command *InCommand) createDirectory(destDir string) error {
	return os.MkdirAll(destDir, 0755)
}

func (command *InCommand) writeURLFile(bucketName string, remotePath string, private bool, destDir string) error {
	url := command.s3client.URL(bucketName, remotePath, private)
	err := ioutil.WriteFile(filepath.Join(destDir, "url"), []byte(url), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (command *InCommand) downloadFile(bucketName string, remotePath string, destinationDir string) error {
	remoteFilename := path.Base(remotePath)
	localPath := filepath.Join(destinationDir, remoteFilename)

	return command.s3client.DownloadFile(
		bucketName,
		remotePath,
		localPath,
	)
}
