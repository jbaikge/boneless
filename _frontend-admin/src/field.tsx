export interface FieldProps {
  type: string;
  name: string;
  label: string;
  sort: boolean;
  column: number;
  min: string;
  max: string;
  step: string;
  format: string;
  options: string;
  class_id: string;
  field: string;
};

export const FieldChoices = [
  { id: 'date',               name: 'Date' },
  { id: 'datetime',           name: 'Date & Time' },
  { id: 'tiny',               name: 'Editor' },
  { id: 'email',              name: 'Email' },
  { id: 'number',             name: 'Number' },
  { id: 'select-class',       name: 'Select (Class)' },
  { id: 'multi-class',        name: 'Mutli-Select (Class)' },
  { id: 'multi-select-label', name: 'Multi-Select w/ Label' },
  { id: 'select-static',      name: 'Select (Static)' },
  { id: 'text',               name: 'Text' },
  { id: 'textarea',           name: 'Textarea' },
  { id: 'time',               name: 'Time' },
  { id: 'any-upload',         name: 'Upload (Any)' },
  { id: 'image-upload',       name: 'Upload (Image)' },
];
