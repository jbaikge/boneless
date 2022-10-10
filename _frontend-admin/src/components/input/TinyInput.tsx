import { useInput, FieldTitle } from 'ra-core';
import { CommonInputProps } from 'react-admin';
import { Editor } from '@tinymce/tinymce-react';

// This implementation may change if there are complaints about performance
// later on. There is a way to update the content when "dirty" instead of every
// keystroke as the current implementation does.
// See: https://www.tiny.cloud/docs/tinymce/6/react-ref/#using-the-tinymce-react-component-as-a-uncontrolled-component
const TinyInput = (props: CommonInputProps) => {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

  const {
    defaultValue = '',
    format,
    parse,
    resource,
    source,
    validate,
    onBlur,
    onChange,
    label,
    helperText,
    ...rest
  } = props;

  const {
    field,
  } = useInput({
    defaultValue,
    format,
    parse,
    resource,
    source,
    type: 'text',
    validate,
    onBlur,
    onChange,
    ...rest,
  });

  return (
    <div style={{ marginBottom: '1em', width: '100%' }}>
      <FieldTitle label={label} source={source} resource={resource} />
      <Editor
        apiKey={process.env.REACT_APP_TINYMCE_KEY}
        init={{
          content_css: (prefersDark ? 'dark' : ''),
          skin: (prefersDark ? 'oxide-dark' : ''),
          width: '100%',
          automatic_uploads: true,
          images_reuse_filename: true,
          images_upload_url: process.env.REACT_APP_API_URL + '/files',
        }}
        plugins={[ 'code', 'image' ]}
        value={field.value}
        onEditorChange={(newValue) => field.onChange(newValue)}
      />
    </div>
  )
};

export default TinyInput;
