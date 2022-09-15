import React from 'react';
import {
  ArrayInput,
  BooleanInput,
  Create,
  CreateProps,
  Datagrid,
  DateField,
  DateInput,
  DateTimeInput,
  Edit,
  EditButton,
  EditProps,
  FormDataConsumer,
  FormDataConsumerRenderParams,
  List,
  ListProps,
  NumberInput,
  ReferenceInput,
  SelectInput,
  SimpleForm,
  SimpleFormIterator,
  TextField,
  TextInput,
  TransformData,
  required,
  useRedirect,
  regex,
  Title,
  FileInput,
  FileField,
  useDataProvider,
  CreateResult,
} from 'react-admin';
import { FieldValues } from 'react-hook-form';
import { Card, CardContent } from '@mui/material';
import { FieldChoices, FieldProps } from './field';
import './App.css';

interface UpdateProps {
  update: React.Dispatch<React.SetStateAction<number>>;
}

interface CreateUpdateProps extends CreateProps, UpdateProps {};

interface EditUpdateProps extends EditProps, UpdateProps {};

interface ClassProps {
  id: string;
  parent_id: string;
  name: string;
  created: string;
  updated: string;
  fields: Array<FieldProps>;
};

const exporter = (classes: Array<ClassProps>) => {
  const blob = new Blob([JSON.stringify(classes, null, 2)], { type: 'application/json' });
  const link = document.createElement('a');
  link.href = window.URL.createObjectURL(blob);
  link.download = 'classes.json';
  link.click();
};

export const ClassCreate = (props: CreateUpdateProps) => {
  const { update, ...rest } = props;
  const redirect = useRedirect();
  const onSuccess = () => {
    update((new Date()).getTime());
    redirect('list', 'classes');
  };

  const ensureFields = (data: TransformData) => ({
    ...data,
    fields: [],
  });

  return (
    <Create {...rest} mutationOptions={{ onSuccess }} transform={ensureFields}>
      <SimpleForm>
        <TextInput source="name" validate={[required()]} fullWidth />
      </SimpleForm>
    </Create>
  );
};

export const ClassEdit = (props: EditUpdateProps) => {
  const { update, ...rest } = props;
  const redirect = useRedirect();
  const onSuccess = () => {
    update((new Date()).getTime());
    redirect('list', 'classes');
  };

  return (
    <Edit {...rest} mutationOptions={{ onSuccess }} mutationMode="pessimistic">
      <SimpleForm>
        <TextInput source="name" validate={[required()]} fullWidth />
        <ReferenceInput source="parent_id" reference="classes" >
          <SelectInput optionText="name" fullWidth />
        </ReferenceInput>
        <ArrayInput source="fields">
          <SimpleFormIterator className="field-row">
            <TextInput source="label" validate={[ required('A label is required') ]} />
            <TextInput source="name" validate={[ required('A field name is required'), regex(/^[a-z0-9_]+$/, 'Names can only contain lowercase letters, numbers and underscores') ]} />
            <BooleanInput source="sort" />
            <NumberInput source="column" />
            <SelectInput source="type" choices={FieldChoices} defaultValue="text" validate={[ required('A type is required') ]} />
            <FormDataConsumer>
              {({
                scopedFormData,
                getSource,
              }: FormDataConsumerRenderParams) => {
                const getSrc = getSource || ((s: string) => s);
                switch (scopedFormData.type) {
                case 'date':
                  return (
                    <>
                      <DateInput source={getSrc('min')} />
                      <DateInput source={getSrc('max')} />
                      <TextInput source={getSrc('step')} label="Step (days)" />
                      <TextInput source={getSrc('format')} label="Format (Jan 2, 2006 3:04pm)" />
                    </>
                  );
                case 'datetime':
                  return (
                    <>
                      <DateTimeInput source={getSrc('min')} />
                      <DateTimeInput source={getSrc('max')} />
                      <TextInput source={getSrc('step')} label="Step (days)" />
                      <TextInput source={getSrc('format')} label="Format (Jan 2, 2006 3:04pm)" />
                    </>
                  );
                case 'time':
                  return (
                    <>
                      <TextInput source={getSrc('format')} label="Format (3:04pm)" />
                    </>
                  );
                case 'number':
                  return (
                    <>
                      <TextInput source={getSrc('min')} />
                      <TextInput source={getSrc('max')} />
                      <TextInput source={getSrc('step')} />
                    </>
                  );
                case 'select-static':
                  return (
                    <>
                      <TextInput source={getSrc('options')} label="Options (one per line, key | value or just value" multiline />
                    </>
                  );
                case 'select-class':
                case 'multi-class':
                case 'multi-select-label':
                  return (
                    <>
                      <ReferenceInput source={getSrc('class_id')} reference="classes">
                        <SelectInput optionText="name" />
                      </ReferenceInput>
                      <TextInput source={getSrc('field')} />
                    </>
                  );
                default:
                  return null;
                }
              }}
            </FormDataConsumer>
          </SimpleFormIterator>
        </ArrayInput>
      </SimpleForm>
    </Edit>
  );
};

export const ClassList = (props: ListProps) => (
  <List {...props} exporter={exporter}>
    <Datagrid>
      <TextField source="name" />
      <DateField source="created" />
      <DateField source="updated" />
      <EditButton />
    </Datagrid>
  </List>
);

export const ClassImport = (props: UpdateProps) => {
  const dataProvider = useDataProvider();
  const redirect = useRedirect();
  let classes: Array<ClassProps> = [];
  const idMap = new Map<string, string>();
  const onSubmit = (data: FieldValues) => {
    const file: File = data['file'].rawFile;
    // Read file data
    file.text()
      // Wait for file to read in as text
      .then((text: string) => JSON.parse(text))
      // Wait for JSON parsing into hopefully the right format.
      // Capture current IDs in idMap.
      .then((data: Array<ClassProps>) => {
        classes = data;
        data.map(item => idMap.set(item.id, ''));
      })
      // Create classes that have no parent_id. Associate new
      // ID with old ID in idMap
      .then(() => Promise.all(
        classes.filter((item: ClassProps) => item.parent_id === '')
        .map((item: ClassProps) => dataProvider.create('classes', { data: { ...item, id: '' } })
          .then((result: CreateResult) => idMap.set(item.id, result.data.id))
        )
      ))
      // Create classes with parent_ids, resetting old ID to new ID
      .then(() => Promise.all(
        classes.filter((item: ClassProps) => item.parent_id !== '')
        .map((item: ClassProps) => {
          const params = {
            data: {
              ...item,
              id: '',
              parent_id: idMap.get(item.parent_id),
            },
          };
          return dataProvider.create('classes', params)
        })
      ))
      // Pull class list and update side navigation
      .then(() => props.update((new Date()).getTime()))
      // We made it! Back to the list page.
      .then(() => redirect('list', 'classes'));
  };

  return (
    <Card sx={{ marginTop: '1em' }}>
      <Title title="Class Import" />
      <CardContent>
        <SimpleForm onSubmit={onSubmit}>
          <FileInput source="file" accept="application/json">
            <FileField source="src" title="title" />
          </FileInput>
        </SimpleForm>
      </CardContent>
    </Card>
  )
}
