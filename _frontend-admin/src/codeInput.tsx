import React from 'react';
import { useInput } from 'ra-core';
import { CommonInputProps } from 'react-admin';
import Editor from 'react-simple-code-editor';
// import { highlight, languages } from 'prismjs';
import * as Prism from 'prismjs';
import 'prismjs/components/prism-markup';
import 'prismjs/components/prism-css';
import 'prismjs/components/prism-clike';
import 'prismjs/components/prism-javascript';

const LightTheme = React.lazy(() => import('./codeInputLight'));
const DarkTheme = React.lazy(() => import('./codeInputDark'));

export const CodeInput = (props: CommonInputProps) => {
  // This mess ganked from TextInput. Would love to know if there is a better
  // way to accomplish this.
  const {
    defaultValue = '',
    label,
    format,
    helperText,
    onBlur,
    onChange,
    parse,
    resource,
    source,
    validate,
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
    <ThemeSelector>
      <Editor
        value={field.value}
        onValueChange={field.onChange}
        highlight={code => Prism.highlight(code, Prism.languages.markup, 'markup')}
        padding={10}
        style={{
          fontFamily: '"JetBrains Mono", "Fira code", "Fira Mono", "monospace"',
          fontSize: 12,
          width: '100%',
          border: '1px solid #888',
        }}
      />
    </ThemeSelector>
  )
}

// Technique from Prawira G
// https://prawira.medium.com/react-conditional-import-conditional-css-import-110cc58e0da6
const ThemeSelector = ({ children }: { children: React.ReactNode }) => {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  return (
    <>
      <React.Suspense fallback={<></>}>
        {!prefersDark && <LightTheme />}
        {prefersDark && <DarkTheme />}
      </React.Suspense>
      {children}
    </>
  )
}

