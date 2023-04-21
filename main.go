package main

import (
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/google/uuid"
)

type Location struct {
	LocationId string  `json:"locationId"`
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Country    string  `json:"country"`
	State      string  `json:"state"`
	City       string  `json:"city"`
	PostalCode string  `json:"postalCode"`
	Address1   string  `json:"address1"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type Project struct {
	ProjectId   string    `json:"projectId"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Team        string    `json:"team"`
	IsComplete  bool      `json:"isComplete"`
	Location    *Location `json:"location"`
	CreatedBy   string    `json:"createdBy"`
	CreatedOn   string    `json:"createdOn"`
	DueDate     string    `json:"dueDate"`
	ModifiedBy  string    `json:"modifiedBy"`
	ModifiedOn  string    `json:"modifiedOn"`
}

// 如果要执行查询，首先要在 CouchBase-UI -> Query Editor 中执行创建索引的语句
// CREATE INDEX idx_projects_team on projects(team);
// CREATE INDEX idx_projects_type on projects(type);
// CREATE INDEX idx_projects_projectId on projects(projectId);

func main() {
	// Uncomment following line to enable logging
	//gocb.SetLogger(gocb.VerboseStdioLogger())

	// Update this to your cluster details
	connectionString := "couchbase://192.168.31.11:30696"
	bucketName := "projects"
	username := "projects-user1"
	password := "user111"
	//username := "Administrator"
	//password := "password"

	// For a secure cluster connection, use `couchbases://<your-cluster-ip>` instead.
	cluster, err := gocb.Connect(connectionString, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
		//TimeoutsConfig: gocb.TimeoutsConfig{
		//	ConnectTimeout: 95 * time.Second,
		//	QueryTimeout:   95 * time.Second,
		//	SearchTimeout:  95 * time.Second,
		//},
	})
	if err != nil {
		log.Fatal(err)
	}

	bucket := cluster.Bucket(bucketName)

	err = bucket.WaitUntilReady(2*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	//Get a reference to the default collection, required for older Couchbase server versions
	col := bucket.DefaultCollection()

	//col := bucket.Scope("_default").Collection("_default")

	// Create and store a Document
	beijingUuid := uuid.New().String()
	_, err = col.Upsert(beijingUuid,
		&Project{
			ProjectId:   "04464b94-c9c8-4726-a941-69c0a62b8dbe",
			Name:        "Warehouse 499",
			Type:        "project",
			Description: "n the other hand, we denounce with righteous indignation and dislike men who are so beguiled and demoralized by the charms of pleasure of the moment, so blinded by desire, that they cannot foresee the pain and trouble that are bound to ensue; and equal blame belongs to those who fail in their duty through weakness of will, which is the same as saying through shrinking from toil and pain. These cases are perfectly simple and easy to distinguish. In a free hour, when our power of choice is untrammelled and when nothing prevents our being able to do what we like best, every pleasure is to be welcomed and every pain avoided. But in certain circumstances and owing to the claims of duty or the obligations of business it will frequently occur that pleasures have to be repudiated and annoyances accepted. The wise man therefore always holds in these matters to this principle of selection: he rejects pleasures to secure other greater pleasures, or else he endures pains to avoid worse pains.",
			Team:        "team4",
			IsComplete:  false,
			Location: &Location{
				LocationId: "a5ae1719-f9d7-4de4-8f44-2ea9d911854d",
				Type:       "location",
				Name:       "Beijing",
				Country:    "China",
				State:      "Beijing",
				City:       "Beijing",
				PostalCode: "100084",
				Address1:   "Zhongguancun",
				Latitude:   40.263834145060547,
				Longitude:  116.6034109424959,
			},
			CreatedBy:  "demo5@example.com",
			CreatedOn:  "1651533549000",
			DueDate:    "1651533549000",
			ModifiedBy: "demo5@example.com",
			ModifiedOn: "1651533549000",
		}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get the document back
	getResult, err := col.Get(beijingUuid, nil)
	if err != nil {
		log.Fatal(err)
	}

	var inProject Project
	err = getResult.Content(&inProject)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Project: %v\n", inProject)

	// Perform a N1QL Query
	//inventoryScope := bucket.Scope("_default")
	query := "SELECT x.* FROM `projects`._default._default x WHERE team=$team and name=$project_name"
	params := make(map[string]interface{}, 1)
	params["team"] = "team4"
	params["project_name"] = "Warehouse 499"
	rows, err := cluster.Query(
		query,
		&gocb.QueryOptions{NamedParameters: params, Adhoc: true},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Print each found Row
	for rows.Next() {
		var result Project
		err := rows.Row(&result)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(result)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
