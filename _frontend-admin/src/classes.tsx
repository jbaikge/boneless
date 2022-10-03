import React from 'react';
import {
  ArrayInput,
  BooleanInput,
  Button,
  Create,
  CreateButton,
  CreateProps,
  CreateResult,
  Datagrid,
  DateField,
  DateInput,
  DateTimeInput,
  Edit,
  EditButton,
  EditProps,
  ExportButton,
  FileField,
  FileInput,
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
  Title,
  TopToolbar,
  TransformData,
  regex,
  required,
  useRedirect,
  useDataProvider,
} from 'react-admin';
import { FieldValues } from 'react-hook-form';
import { Link } from 'react-router-dom';
import IconFileUpload from '@mui/icons-material/FileUpload';
import { Card, CardContent } from '@mui/material';
import { jsonExporter } from './exporter';
import { FieldChoices, FieldProps } from './field';
import { GlobalPagination } from './pagination';
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

export const ClassCreate = (props: CreateUpdateProps) => {
  const { update, ...rest } = props;
  const redirect = useRedirect();
  const onSuccess = (data: ClassProps) => {
    update((new Date()).getTime());
    redirect('edit', 'classes', data.id);
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
            <BooleanInput source="sort" defaultValue={false} label="Index this data for sorting" />
            <NumberInput source="column" defaultValue={0} label="List Column (0 = hidden)" />
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
                      <TextInput source={getSrc('format')} defaultValue="Jan 2, 2006" label="Format (Jan 2, 2006)" />
                    </>
                  );
                case 'datetime':
                  return (
                    <>
                      <DateTimeInput source={getSrc('min')} />
                      <DateTimeInput source={getSrc('max')} />
                      <TextInput source={getSrc('step')} label="Step (days)" />
                      <TextInput source={getSrc('format')} defaultValue="Jan 2, 2006 3:04pm" label="Format (Jan 2, 2006 3:04pm)" />
                    </>
                  );
                case 'time':
                  return (
                    <>
                      <TextInput source={getSrc('format')} defaultValue="3:04pm" label="Format (3:04pm)" />
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

const ListActions = () => (
  <TopToolbar>
    <CreateButton />
    <ExportButton />
    <Button label="Import" component={Link} to="/class-import">
      <IconFileUpload />
    </Button>
  </TopToolbar>
);

export const ClassList = (props: ListProps) => (
  <List {...props} actions={<ListActions />} exporter={jsonExporter('classes')} pagination={<GlobalPagination />}>
    <Datagrid sx={{
      '& td:nth-last-of-type(2)': { width: '8em' },
      '& td:last-child': { width: '5em' },
    }}>
      <TextField source="name" />
      <DateField source="created" />
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
