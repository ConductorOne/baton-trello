package client

type PaginationVars struct {
	Size uint
	Page uint
}

type User struct {
	ID         string `json:"id"`
	MemberID   string `json:"idMember"`
	Name       string `json:"fullName"`
	Username   string `json:"username"`
	MemberType string `json:"memberType"`
}

type Organization struct {
	ID              string `json:"id"`
	DisplayName     string `json:"displayName"`
	Name            string `json:"name"`
	Description     string `json:"desc"`
	DescriptionData struct {
		Emoji struct{} `json:"emoji"`
	} `json:"descData"`
	Url      string `json:"url"`
	Website  string `json:"website"`
	TeamType string `json:"teamType"`
	LogoHash string `json:"logoHash"`
	LogoUrl  string `json:"logoUrl"`
	Offering string `json:"offering"`
	Products []int  `json:"products"`
	PowerUps []int  `json:"powerUps"`
}
