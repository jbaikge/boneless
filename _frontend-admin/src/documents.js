import {
    Create,
    Datagrid,
    Edit,
    List
} from 'react-admin';

export const DocumentCreate = (props) => (
    <Create {...props}>
    </Create>
);

export const DocumentEdit = (props) => (
    <Edit {...props}>
    </Edit>
);

export const DocumentList = (props) => (
    <List {...props}>
        <Datagrid rowClick="edit">
        </Datagrid>
    </List>
);
