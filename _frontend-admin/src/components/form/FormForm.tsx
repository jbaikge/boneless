import {
  ArrayInput,
  BooleanInput,
  FormTab,
  NumberInput,
  SelectInput,
  SimpleFormIterator,
  TabbedForm,
  TextInput,
  required,
} from 'react-admin';

const FormForm = () => {
  return (
    <TabbedForm>
      <FormTab label="Settings">
        <TextInput source="id" label="Form ID" helperText="Available after creation" disabled fullWidth />
        <TextInput source="name" fullWidth />
      </FormTab>
      <FormTab label="Fields">
        <ArrayInput source="fields">
          <SimpleFormIterator className="field-row">
            <TextInput source="label" validate={[ required('A label is required') ]} />
            <TextInput source="name" validate={[ required('A name is required') ]} />
            <BooleanInput source="required" defaultValue={false} label="Required" />
            <NumberInput source="column" defaultValue={0} label="List Column (0 = hidden)" />
            <SelectInput source="type" choices={[{ id: 'text', name: 'Text' }]} defaultValue="text" validate={[ required('A type is required') ]} />
            <ArrayInput source="validations">
              <SimpleFormIterator>
                <SelectInput source="type" choices={[{id: 'required', name: 'Required'}]}></SelectInput>
              </SimpleFormIterator>
            </ArrayInput>
          </SimpleFormIterator>
        </ArrayInput>
      </FormTab>
    </TabbedForm>
  );
};

export default FormForm;
