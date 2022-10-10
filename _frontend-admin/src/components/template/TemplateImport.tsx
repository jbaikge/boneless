import {
  CreateResult,
  FileField,
  FileInput,
  SimpleForm,
  Title,
  useDataProvider,
  useRedirect,
} from 'react-admin';
import { FieldValues } from 'react-hook-form';
import { Card, CardContent } from '@mui/material';
import { TemplateProps } from './Props';

const TemplateImport = () => {
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
      <Title title="Template Import" />
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

export default TemplateImport;
