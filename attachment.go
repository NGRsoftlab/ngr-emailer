// Copyright 2020-2023 NGR Softlab
//
package emailer

// AttachData struct for attachments
type AttachData struct {
	FileName string
	FileData []byte
}

// attachFile attaching file to mail-obj. Returns attachments map: "filename": filedata
func attachFile(files []AttachData) (map[string][]byte, error) {
	var attachments = make(map[string][]byte)
	for _, f := range files {
		attachments[f.FileName] = f.FileData
	}
	return attachments, nil
}
