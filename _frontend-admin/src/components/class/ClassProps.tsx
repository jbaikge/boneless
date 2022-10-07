import FieldProps from '../field/FieldProps';

interface ClassProps {
  id: string;
  parent_id: string;
  name: string;
  created: string;
  updated: string;
  fields: Array<FieldProps>;
};

export default ClassProps;
