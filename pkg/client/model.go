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

type Preferences struct {
	PermissionLevel          string        `json:"permissionLevel"`
	HideVotes                bool          `json:"hideVotes"`
	Voting                   string        `json:"voting"`
	Comments                 string        `json:"comments"`
	Invitations              string        `json:"invitations"`
	SelfJoin                 bool          `json:"selfJoin"`
	CardCovers               bool          `json:"cardCovers"`
	ShowCompleteStatus       bool          `json:"showCompleteStatus"`
	CardCounts               bool          `json:"cardCounts"`
	IsTemplate               bool          `json:"isTemplate"`
	CardAging                string        `json:"cardAging"`
	CalendarFeedEnabled      bool          `json:"calendarFeedEnabled"`
	HiddenPluginBoardButtons []interface{} `json:"hiddenPluginBoardButtons"`
	SwitcherViews            []struct {
		ViewType string `json:"viewType"`
		Enabled  bool   `json:"enabled"`
	} `json:"switcherViews"`
	Background            string      `json:"background"`
	BackgroundColor       string      `json:"backgroundColor"`
	BackgroundDarkColor   interface{} `json:"backgroundDarkColor"`
	BackgroundImage       interface{} `json:"backgroundImage"`
	BackgroundDarkImage   interface{} `json:"backgroundDarkImage"`
	BackgroundImageScaled interface{} `json:"backgroundImageScaled"`
	BackgroundTile        bool        `json:"backgroundTile"`
	BackgroundBrightness  string      `json:"backgroundBrightness"`
	SharedSourceUrl       interface{} `json:"sharedSourceUrl"`
	BackgroundBottomColor string      `json:"backgroundBottomColor"`
	BackgroundTopColor    string      `json:"backgroundTopColor"`
	CanBePublic           bool        `json:"canBePublic"`
	CanBeEnterprise       bool        `json:"canBeEnterprise"`
	CanBeOrg              bool        `json:"canBeOrg"`
	CanBePrivate          bool        `json:"canBePrivate"`
	CanInvite             bool        `json:"canInvite"`
}

type Board struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Description    string      `json:"desc"`
	DescData       interface{} `json:"descData"`
	Closed         bool        `json:"closed"`
	IdOrganization string      `json:"idOrganization"`
	IdEnterprise   interface{} `json:"idEnterprise"`
	Pinned         bool        `json:"pinned"`
	Url            string      `json:"url"`
	ShortUrl       string      `json:"shortUrl"`
	Preferences    Preferences `json:"prefs"`
	LabelNames     struct {
		Green       string `json:"green"`
		Yellow      string `json:"yellow"`
		Orange      string `json:"orange"`
		Red         string `json:"red"`
		Purple      string `json:"purple"`
		Blue        string `json:"blue"`
		Sky         string `json:"sky"`
		Lime        string `json:"lime"`
		Pink        string `json:"pink"`
		Black       string `json:"black"`
		GreenDark   string `json:"green_dark"`
		YellowDark  string `json:"yellow_dark"`
		OrangeDark  string `json:"orange_dark"`
		RedDark     string `json:"red_dark"`
		PurpleDark  string `json:"purple_dark"`
		BlueDark    string `json:"blue_dark"`
		SkyDark     string `json:"sky_dark"`
		LimeDark    string `json:"lime_dark"`
		PinkDark    string `json:"pink_dark"`
		BlackDark   string `json:"black_dark"`
		GreenLight  string `json:"green_light"`
		YellowLight string `json:"yellow_light"`
		OrangeLight string `json:"orange_light"`
		RedLight    string `json:"red_light"`
		PurpleLight string `json:"purple_light"`
		BlueLight   string `json:"blue_light"`
		SkyLight    string `json:"sky_light"`
		LimeLight   string `json:"lime_light"`
		PinkLight   string `json:"pink_light"`
		BlackLight  string `json:"black_light"`
	} `json:"labelNames"`
	Memberships []User `json:"memberships"`
}
