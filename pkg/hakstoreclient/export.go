package hakstoreclient

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

// Export will export the whole database to markdown files for obsidian
func ExportCLI(c Client) {
	exportFlagSet := flag.NewFlagSet("export", flag.ExitOnError)
	outputDirPtr := exportFlagSet.String("d", ".", "directory to store output files")
	outputDir := *outputDirPtr
	exportFlagSet.Parse(os.Args[2:])

	err := os.Mkdir(outputDir, 0755)
	if err != nil {
		log.Println("Error creating directory:", err)
	}
	platforms, err := c.GetPlatforms()
	if err != nil {
		log.Println("Error encountered retrieving platforms:", err)
	}
	for _, platform := range platforms {
		err := os.Mkdir(outputDir+"/"+platform.ID, 0755)
		if err != nil {
			log.Println("Error creating directory:", err)
		}
		fileContents := []byte("# " + platform.ID)
		err = ioutil.WriteFile(outputDir+"/"+platform.ID+"/"+platform.ID+".md", fileContents, 0644)
		if err != nil {
			log.Println("Error writing file:", err)
		}
		programs, err := c.GetAssociatedPrograms(platform.ID)
		if err != nil {
			log.Println("Error retrieving programs", err)
		}
		for _, program := range programs {
			err := os.Mkdir(outputDir+"/"+platform.ID+"/"+program.ID, 0755)
			if err != nil {
				log.Println("Error creating directory:", err)
			}
			fileContents := []byte("# " + program.ID + "\n\nPlatform: [[" + program.PlatformID + "]]")
			err = ioutil.WriteFile(outputDir+"/"+platform.ID+"/"+program.ID+"/"+program.ID+".md", fileContents, 0644)
			if err != nil {
				log.Println("Error writing file:", err)
			}
			rootdomains, err := c.GetAssociatedRootDomains(program.ID)
			if err != nil {
				log.Println("Error retrieving programs", err)
			}
			for _, rootdomain := range rootdomains {
				err := os.Mkdir(outputDir+"/"+platform.ID+"/"+program.ID+"/"+rootdomain.ID, 0755)
				if err != nil {
					log.Println("Error creating directory:", err)
				}
				fileContents := []byte("# " + program.ID + "\n\nPlatform: [[" + platform.ID + "]]\nProgram: [[" + program.ID + "]]")
				err = ioutil.WriteFile(outputDir+"/"+platform.ID+"/"+program.ID+"/"+rootdomain.ID+"/"+rootdomain.ID+".md", fileContents, 0644)
				if err != nil {
					log.Println("Error writing file:", err)
				}
				subdomains, err := c.GetAssociatedSubdomains(rootdomain.ID)
				if err != nil {
					log.Println("Error retrieving programs", err)
				}
				for _, subdomain := range subdomains {
					fileContents := []byte("# " + program.ID + "\n\nPlatform: [[" + platform.ID + "]]\nProgram: [[" + program.ID + "]]\nRoot Domain: [[" + rootdomain.ID + "]]")
					err = ioutil.WriteFile(outputDir+"/"+platform.ID+"/"+program.ID+"/"+rootdomain.ID+"/"+subdomain.ID+".md", fileContents, 0644)
					if err != nil {
						log.Println("Error writing file:", err)
					}
				}
			}
		}
	}
}
