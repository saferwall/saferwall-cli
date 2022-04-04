// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package cmd

import (
	"io/ioutil"
	"log"

	"github.com/saferwall/saferwall-cli/internal/entity"
	"github.com/saferwall/saferwall-cli/internal/webapi"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Used for flags.
var corpusYaml string

func init() {
	genCmd.Flags().StringVarP(&corpusYaml, "corpus", "c", "",
		"Yaml corpus file (required)")
	genCmd.MarkFlagRequired("corpus")
}

func loadCorpus(filename string) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read corpus yaml file :%v ", err)
	}

	var corpus entity.Corpus
	err = yaml.Unmarshal(yamlFile, &corpus)
	if err != nil {
		log.Fatalf("failed to unmarshal yaml string: %v", err)
	}

	for _, fam := range corpus.Families {
		log.Printf("processing %s", fam.Name)
		for _, sample := range fam.Samples {
			var file entity.File

			log.Printf("processing %s | %s | %s | %s",
				sample.SHA256, sample.Platform, sample.FileFormat, sample.Category)

			err = webapi.GetFile(sample.SHA256, &file)
			if err != nil {
				log.Fatalf("failed to read doc from saferwall web service: %v", err)
			}
		}
	}
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate malware souk markdown for the entire corpus",
	Long: `Generates markdown source code for the entire corpus of
saferwall's malware souk database`,
	Run: func(cmd *cobra.Command, args []string) {
		loadCorpus(corpusYaml)
	},
}
