import { RichTextInput } from 'ra-input-rich-text';
import React from 'react';
import {
  ArrayField,
  ArrayInput,
  Create,
  CreateProps,
  Datagrid,
  DateField,
  DateTimeInput,
  Edit,
  EditButton,
  EditProps,
  FileField,
  FileInput,
  ImageField,
  ImageInput,
  List,
  ListProps,
  Loading,
  ReferenceField,
  ReferenceInput,
  ReferenceManyField,
  RichTextField,
  SelectInput,
  Show,
  ShowButton,
  ShowProps,
  SimpleForm,
  SimpleFormIterator,
  SimpleShowLayout,
  TextField,
  TextInput,
  useGetList,
  useGetOne,
  useRecordContext,
  useResourceContext,
} from 'react-admin';
import { Box, Typography } from '@mui/material';
import { FieldProps } from './field';
import { GlobalPagination } from './pagination';
import { TinyInput } from './tinyInput';

const resourceRE = /([^/]+)\/([^/]+)\/.*/g;

export const DocumentCreate = (props: CreateProps) => {
  return (
    <Create {...props}>
      <DocumentForm />
    </Create>
  );
};

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

export const DocumentEdit = (props: EditProps) => {
  return (
    <Edit {...props} aside={<EditAside />}>
      <DocumentForm />
    </Edit>
  );
};

export const DocumentForm = () => {
  const resourceContext = useResourceContext();
  // resourceContext should be "classes/<id>/documents"
  const [[ , resource, id ]] = [...resourceContext.matchAll(resourceRE)];
  const { data, isLoading } = useGetOne(resource, { id });

  if (isLoading) {
    return <Loading />;
  }

  let parentInput;
  if (data.parent_id !== "") {
    parentInput = <ReferenceInput source="parent_id" reference={'classes/' + data.parent_id + '/documents'} perPage={25}>
      <SelectInput optionText="values.title" fullWidth />
    </ReferenceInput>;
  }

  return (
    <SimpleForm>
      <TextInput source="path" fullWidth />
      {parentInput}
      <ReferenceInput source="template_id" reference="templates" perPage={100}>
        <SelectInput optionText="name" fullWidth />
      </ReferenceInput>
      {data.fields.map((field: FieldProps) => {
        const source = `values.${field.name}`;
        switch (field.type) {
          case 'datetime':
            return <DateTimeInput key={field.name} source={source} label={field.label} inputProps={{ min: field.min, max: field.max, step: field.step }} />
          case 'multi-class':
            return (
              <ArrayInput key={field.name} source={source} label={field.label}>
                <SimpleFormIterator>
                  <ReferenceInput reference={'/classes/' + field.class_id + '/documents'} source="id" perPage={100} sort={{ field: field.field, order: 'ASC' }}>
                    <SelectInput optionText={'values.' + field.field} label={field.label} />
                  </ReferenceInput>
                </SimpleFormIterator>
              </ArrayInput>
            );
          case 'multi-select-label':
            return (
              <ArrayInput source={source} label={field.label}>
                <SimpleFormIterator>
                  <ReferenceInput reference={'/classes/' + field.class_id + '/documents'} source='id' perPage={100} sort={{ field: field.field, order: 'ASC' }}>
                    <SelectInput optionText={'values.' + field.field} label={field.label} />
                  </ReferenceInput>
                  <TextInput source='label' />
                </SimpleFormIterator>
              </ArrayInput>
            );
          case 'richtext':
            return <RichTextInput key={field.name} source={source} label={field.label} fullWidth />
          case 'select-class':
            return (
              <ReferenceInput reference={'/classes/' + field.class_id + '/documents'} source={source} perPage={100} sort={{ field: field.field, order: 'ASC' }}>
                <SelectInput optionText={'values.' + field.field} label={field.label} fullWidth />
              </ReferenceInput>
            );
          case 'text':
            return <TextInput key={field.name} source={source} label={field.label} fullWidth />
          case 'textarea':
            return <TextInput key={field.name} source={source} label={field.label} fullWidth multiline />
          case 'time':
            return <TextInput type="time" key={field.name} source={source} label={field.label} fullWidth />
          case 'tiny':
            return <TinyInput key={field.name} source={source} label={field.label} fullWidth />
          case 'any-upload':
            return (
              <FileInput key={field.name} source={source} label={field.label} fullWidth>
                <FileField source="url" title="title" />
              </FileInput>
            );
          case 'image-upload':
            return (
              <React.Fragment key={field.name}>
                <ImageInput source={source} label={field.label} fullWidth>
                    <ImageField source="url" title="title" />
                </ImageInput>
                <TextInput source={source + '.path'} label={field.label + ' Path (start with a slash)'} fullWidth />
              </React.Fragment>
            );
          default:
            return <TextInput key={field.name} source={source} label={`Unknown type (${field.type}) - ${field.label}`} fullWidth />
        }
      })}
    </SimpleForm>
  );
};

export const DocumentList = (props: ListProps) => {
  const resourceContext = useResourceContext();
  const [[ , resource, id ]] = [...resourceContext.matchAll(resourceRE)];
  const { data, isLoading } = useGetOne(resource, { id });

  if (isLoading) {
    return <Loading />;
  }

  let parentField = null;
  if (data.parent_id !== '') {
    parentField = (
      <ReferenceField reference={`classes/${data.parent_id}/documents`} source="parent_id" label="Parent" sortable={false}>
        <TextField source="values.title" />
      </ReferenceField>
    );
  }

  return (
    <List {...props} pagination={<GlobalPagination />}>
      <Datagrid sx={{
        '& td:last-child': { width: '5em' },
        '& td:nth-last-child(2)': { width: '5em' },
      }}>
        {parentField}
        {data.fields.filter((field: FieldProps) => field.column > 0).sort((a: FieldProps, b: FieldProps) => a.column - b.column).map((field: FieldProps) => {
          const source = `values.${field.name}`;
          switch (field.type) {
            case 'date':
              return <DateField key={field.label} source={source} label={field.label} sortable={field.sort} />
            case 'datetime':
              return <DateField key={field.label} source={source} label={field.label} sortable={field.sort} showTime />
            case 'image-upload':
              return <ImageField key={field.label} source={`${source}.url`} label={field.label} sortable={field.sort} />
            default:
              return <TextField key={field.label} source={source} label={field.label} sortable={field.sort} />
          }
        })}
        <ShowButton />
        <EditButton />
      </Datagrid>
    </List>
  );
};

export const DocumentShow = (props: ShowProps) => {
  const resourceContext = useResourceContext();
  const [[ , resource, id ]] = [...resourceContext.matchAll(resourceRE)];
  const { data, isLoading } = useGetOne(resource, { id });

  if (isLoading) {
    return <Loading />;
  }

  let parentField = null;
  if (data.parent_id !== '') {
    parentField = (
      <ReferenceField reference={`classes/${data.parent_id}/documents`} source="parent_id" label="Parent" sortable={false}>
        <TextField source="values.title" />
      </ReferenceField>
    );
  }

  return (
    <Show {...props}>
      <SimpleShowLayout>
        {parentField}
        {data.fields.map((field: FieldProps) => {
          const source = `values.${field.name}`;
          const label = `${field.label} (.Document.Values.${field.name})`
          switch (field.type) {
            case 'image-upload':
              return <ImageField source={`${source}.url`} label={label} />
            case 'multi-class':
              return <ArrayField source={source} label={false}>
                <Datagrid bulkActionButtons={false}>
                  <ReferenceField reference={`classes/${field.class_id}/documents`} source="id" label={label}>
                    <TextField source={`values.${field.field}`} />
                  </ReferenceField>
                </Datagrid>
              </ArrayField>
            case 'tiny':
              return <RichTextField source={source} label={label} />
            default:
              return <TextField source={source} label={label} />
          }
        })}
      </SimpleShowLayout>
    </Show>
  )
}
