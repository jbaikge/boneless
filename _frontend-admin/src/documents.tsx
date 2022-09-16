import { RichTextInput } from 'ra-input-rich-text';
import React from 'react';
import {
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
  SelectInput,
  Show,
  ShowButton,
  ShowProps,
  SimpleForm,
  SimpleFormIterator,
  SimpleShowLayout,
  TextField,
  TextInput,
  useGetOne,
  useResourceContext,
} from 'react-admin';
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

export const DocumentEdit = (props: EditProps) => {
  return (
    <Edit {...props}>
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
        <SelectInput fullWidth />
      </ReferenceInput>
      {data.fields.map((field: FieldProps) => {
        const source = `values.${field.name}`;
        switch (field.type) {
          case 'datetime':
            return <DateTimeInput key={field.name} source={source} label={field.label} inputProps={{ min: field.min, max: field.max, step: field.step }} />
          case 'multi-class':
            return (
              <ArrayInput source={source} label={field.label}>
                <SimpleFormIterator>
                  <ReferenceInput reference={'/classes/' + field.class_id + '/documents'} source="id">
                    <SelectInput optionText={'values.' + field.field} label={field.label} />
                  </ReferenceInput>
                </SimpleFormIterator>
              </ArrayInput>
            );
          case 'multi-select-label':
            return (
              <ArrayInput source={source} label={field.label}>
                <SimpleFormIterator>
                  <ReferenceInput reference={'/classes/' + field.class_id + '/documents'} source='id'>
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
              <ReferenceInput reference={'/classes/' + field.class_id + '/documents'} source={source}>
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
      <Datagrid>
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
        <EditButton />
        <ShowButton />
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
          switch (field.type) {
            default:
              return <TextField key={field.label} source={source} label={field.label} />
          }
        })}
      </SimpleShowLayout>
    </Show>
  )
}
