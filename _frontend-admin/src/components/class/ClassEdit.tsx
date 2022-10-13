import {
    ArrayInput,
    BooleanInput,
    DateInput,
    DateTimeInput,
    Edit,
    FormDataConsumer,
    FormDataConsumerRenderParams,
    NumberInput,
    ReferenceInput,
    SelectInput,
    SimpleForm,
    SimpleFormIterator,
    TextInput,
    regex,
    required,
    useRedirect,
} from 'react-admin';
import { EditUpdateProps } from './Props';
import { FieldChoices } from '../field';

const ClassEdit = (props: EditUpdateProps) => {
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

export default ClassEdit;
