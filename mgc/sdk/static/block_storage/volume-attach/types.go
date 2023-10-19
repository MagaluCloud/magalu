package attachment

type VolumeAttachmentResponse struct {
	VirtualMachineId string `json:"virtual_machine_id"`
}

type VolumeResponse struct {
	Attachments []VolumeAttachmentResponse `json:"attachments"`
	ID          string                     `json:"id"`
}

type AttachmentResult struct {
	VolumeID         string `json:"id"`
	VirtualMachineID string `json:"virtual_machine_id"`
}
