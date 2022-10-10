import {
  Button,
  CreateButton,
  Datagrid,
  DateField,
  EditButton,
  ExportButton,
  List,
  ListProps,
  TextField,
  TopToolbar,
} from 'react-admin';
import { Link } from 'react-router-dom';
import IconFileUpload from '@mui/icons-material/FileUpload';
import JSONExport from '../export/JSONExport';
import GlobalPagination from '../GlobalPagination';

const ListActions = () => (
  <TopToolbar>
    <CreateButton />
    <ExportButton />
    <Button label="Import" component={Link} to="/template-import">
      <IconFileUpload />
    </Button>
  </TopToolbar>
);

const TemplateList = (props: ListProps) => (
  <List {...props} actions={<ListActions />} exporter={JSONExport('templates')} pagination={<GlobalPagination />}>
    <Datagrid sx={{
      '& td:nth-last-of-type(3)': { width: '8em' },
      '& td:nth-last-of-type(2)': { width: '8em' },
      '& td:last-child': { width: '5em' },
    }}>
      <TextField source="name" />
      <DateField source="created" />
      <DateField source="updated" />
      <EditButton />
    </Datagrid>
  </List>
);

export default TemplateList;
