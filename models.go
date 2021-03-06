package s3resource

type Source struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Bucket          string `json:"bucket"`
	Regexp          string `json:"regexp"`
	Private         bool   `json:"private"`
	RegionName      string `json:"region_name"`
}

type Version struct {
	Path string `json:"path"`
}

type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
