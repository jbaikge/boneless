import {
  Create,
  CreateProps,
  CreateResult,
  Datagrid,
  DateField,
  Edit,
  EditButton,
  EditProps,
  FileField,
  FileInput,
  List,
  ListProps,
  SimpleForm,
  TextField,
  TextInput,
  Title,
  useDataProvider,
  useRedirect,
} from 'react-admin';
import { FieldValues } from 'react-hook-form';
import { Card, CardContent } from '@mui/material';
import { CodeInput } from './codeInput';
import { jsonExporter } from './exporter';
import { GlobalPagination } from './pagination';

interface TemplateProps {
  id: string;
  name: string;
  version: number;
  body: string;
  created: string;
  updated: string;
}

export const TemplateCreate = (props: CreateProps) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="name" fullWidth />
      <CodeInput source="body" />
    </SimpleForm>
  </Create>
);

export const TemplateEdit = (props: EditProps) => (
  <Edit {...props}>
    <SimpleForm>
      <TextInput source="name" fullWidth />
      <CodeInput source="body" />
    </SimpleForm>
  </Edit>
);

export const TemplateList = (props: ListProps) => (
  <List {...props} exporter={jsonExporter('templates')} pagination={<GlobalPagination />}>
    <Datagrid>
      <TextField source="name" />
      <DateField source="created" />
      <DateField source="updated" />
      <EditButton />
    </Datagrid>
  </List>
);

export const TemplateImport = () => {
  const dataProvider = useDataProvider();
  const redirect = useRedirect();
  let templates: Array<TemplateProps> = [];
  const idMap = new Map<string, string>();
  const onSubmit = (data: FieldValues) => {
    const file: File = data['file'].rawFile;
    // Read file data
    file.text()
      // Wait for file to read in as text
      .then((text: string) => JSON.parse(text))
      // Wait for JSON parsing into hopefully the right format.
      // Capture current IDs in idMap.
      .then((data: Array<TemplateProps>) => {
        templates = data;
        data.map(item => idMap.set(item.id, ''));
      })
      // Create templates that have no parent_id. Associate new
      // ID with old ID in idMap
      .then(() => Promise.all(
        templates.map((item: TemplateProps) => dataProvider.create('templates', { data: { ...item, id: '' } })
          .then((result: CreateResult) => idMap.set(item.id, result.data.id))
        )
      ))
      // We made it! Back to the list page.
      .then(() => redirect('list', 'templates'));
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
};
