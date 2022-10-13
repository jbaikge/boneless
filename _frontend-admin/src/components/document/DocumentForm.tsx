import React from 'react';
import {
  ArrayInput,
  DateTimeInput,
  FileField,
  FileInput,
  ImageField,
  ImageInput,
  Loading,
  ReferenceInput,
  SelectInput,
  SimpleForm,
  SimpleFormIterator,
  TextInput,
  useGetOne,
  useResourceContext,
} from 'react-admin';
import { reResource } from './Constants';
import { FieldProps } from '../field/Props';
import { TinyInput } from '../input';

const DocumentForm = () => {
  const resourceContext = useResourceContext();
  // resourceContext should be "classes/<id>/documents"
  const [[ , resource, id ]] = [...resourceContext.matchAll(reResource)];
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

export default DocumentForm;
