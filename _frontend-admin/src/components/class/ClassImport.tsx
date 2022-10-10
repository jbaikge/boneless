import {
  CreateResult,
  FileField,
  FileInput,
  SimpleForm,
  Title,
  useDataProvider,
  useRedirect,
} from 'react-admin';
import { Card, CardContent } from '@mui/material';
import { FieldValues } from 'react-hook-form';
import { ClassProps, UpdateProps } from './Props';

const ClassImport = (props: UpdateProps) => {
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

export default ClassImport;
