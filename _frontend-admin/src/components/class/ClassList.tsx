import {
  Button,
  CreateButton,
  Datagrid,
  EditButton,
  DateField,
  ExportButton,
  List,
  ListProps,
  TextField,
  TopToolbar,
} from 'react-admin';
import { Link } from 'react-router-dom';
import IconFileUpload from '@mui/icons-material/FileUpload';
import GlobalPagination from '../GlobalPagination';

export const jsonExporter = (baseName: string) => (objects: Array<any>) => {
  const blob = new Blob([JSON.stringify(objects, null, 2)], { type: 'application/json' });
  const link = document.createElement('a');
  link.href = window.URL.createObjectURL(blob);
  link.download = `${baseName}.json`;
  link.click();
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

const ClassList = (props: ListProps) => (
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

export default ClassList;
