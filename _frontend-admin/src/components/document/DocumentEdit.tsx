import React from 'react';
import {
  Datagrid,
  Edit,
  EditProps,
  Loading,
  ReferenceManyField,
  TextField,
  useGetList,
  useRecordContext,
} from 'react-admin';
import { Box, Typography } from '@mui/material';
import DocumentForm from './DocumentForm';

const EditAside = () => {
  const record = useRecordContext();
  const { data, isLoading } = useGetList('classes');

  if (isLoading || !record) {
    return <Loading />;
  }

  const filtered = data?.filter((value) => value.parent_id === record.class_id);
  if (filtered?.length === 0) {
    return <></>
  }

  return (
    <Box sx={{ width: '240px', margin: '1em' }}>
      {filtered?.map((c) => {
        const resource = `classes/${c.id}/documents`;
        return (
          <React.Fragment key={c.id}>
            <Typography variant="h6">{c.name}</Typography>
            <ReferenceManyField reference={resource} target="parent_id">
              <Datagrid bulkActionButtons={false} rowClick="edit">
                <TextField source="values.title" label="Title" />
              </Datagrid>
            </ReferenceManyField>
          </React.Fragment>
        );
      })}
    </Box>
  );
}

const DocumentEdit = (props: EditProps) => {
  return (
    <Edit {...props} aside={<EditAside />}>
      <DocumentForm />
    </Edit>
  );
};

export default DocumentEdit;
