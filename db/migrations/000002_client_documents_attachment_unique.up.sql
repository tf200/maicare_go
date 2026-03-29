CREATE UNIQUE INDEX client_documents_attachment_uuid_unique_idx
ON client_documents(attachment_uuid)
WHERE attachment_uuid IS NOT NULL;
