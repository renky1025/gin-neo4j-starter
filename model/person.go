package model

type PersonModel struct {
	Id        int64  `json:"id" from:"id"`
	Number    int64  `json:"number" from:"number"`
	Apperance int64  `json:"apperance" from:"apperance"`
	Position  string `json:"position" from:"position"`
	Name      string `json:"name" from:"name"`
	Age       int64  `json:"age" from:"age"`
	Goal      int64  `json:"goal" from:"goal"`
}

type NodeFrom struct {
	Label string `json:"label" from:"label"`
	Name  string `json:"name" from:"name"`
}

type RelationMap struct {
	RelationShip string                 `json:"relationShip" from:"relationShip"`
	Attr         map[string]interface{} `json:"attr" from:"attr"`
}
