import {
    useListContext,
    CreateButton,
    Datagrid,
    DateField,
    ExportButton,
    List,
    TextField,
    TopToolbar,
    Create,
    SimpleForm,
    TextInput,
    required
} from 'react-admin';

const ListActions = () => {
    const { total, isLoading } = useListContext();

    return (
        <TopToolbar>
            <CreateButton />
            <ExportButton disabled={isLoading || total === 0} />
        </TopToolbar>
    );
}

export const ClassCreate = () => (
    <Create>
        <SimpleForm>
            <TextInput source="name" validate={[required()]} fullWidth />
            <TextInput source="slug" validate={[required()]} fullWidth />
            <TextInput source="table_labels" />
            <TextInput source="table_fields" />
        </SimpleForm>
    </Create>
)

export const ClassList = () => (
    <List>
        <Datagrid rowClick="edit">
            <TextField source="id" />
            <TextField source="name" />
            <TextField source="slug" />
            <DateField source="created" />
            <DateField source="updated" />
        </Datagrid>
    </List>
);