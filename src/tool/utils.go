package tool

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func generate_file_name(migration_name string) string {
	timestamp := int (time.Now().UnixMicro())
	return file_name_format(strconv.Itoa(timestamp), migration_name)
}

func file_name_format(id string, migration_name string) string {
	return fmt.Sprintf("%s_%s.sql", id, migration_name)
}

func parse_file_name(file_name string) (ParsedFileName, error) {
	split_file_name := strings.SplitN(file_name, "_", 2)
	if len(split_file_name) != 2 {
		return ParsedFileName{}, fmt.Errorf("incorrect file format")
	}
	
	// file name format: "{id}_{migration_name}.sql"
	id := split_file_name[0]
	split_file_name_extension := strings.SplitN(split_file_name[1], ".", 2)
	if len(split_file_name_extension) != 2 {
		return ParsedFileName{}, fmt.Errorf("incorrect file format")
	}

	migration_name := split_file_name_extension[0]
	file_extension := split_file_name_extension[1]
	return ParsedFileName{
		Id: id,
		MigrationName: migration_name,
		FileExtension: file_extension,
		Raw: file_name,
	}, nil
}