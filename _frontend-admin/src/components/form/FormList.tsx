import {
  Datagrid,
  EditButton,
  List,
  ListProps,
  ShowButton,
  TextField,
} from 'react-admin';
import GlobalPagination from '../GlobalPagination';

const FormList = (props: ListProps) => {
  return (
    <List {...props} pagination={<GlobalPagination />}>
      <Datagrid sx={{
        '& td:last-child': { width: '5em' },
        '& td:nth-last-of-type(2)': { width: '5em' },
      }}>
        <TextField source="name" label="Form Name" />
        <ShowButton />
        <EditButton />
      </Datagrid>
    </List>
  )
}

export default FormList;
