package main

import (
	"fmt"
	"github.com/Abathargh/moody-go/db"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

func main() {
	if err := db.Init(); err != nil {
		log.Println("Couldn't initialize the database!")
		log.Fatal(err)
	}

	nodes, rowCount, err := db.GetAllNodes(nil, 0, 20, "")

	fmt.Println(rowCount)

	for _, node := range nodes {
		fmt.Println(node)
	}

	mac := "167317910391173"
	node, err := db.GetNode(nil, mac)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(node)

}
