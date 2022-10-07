import React from 'react';
import { CreateProps, EditProps } from 'react-admin';
import { FieldProps } from '../field/Props';

export interface ClassProps {
  id: string;
  parent_id: string;
  name: string;
  created: string;
  updated: string;
  fields: Array<FieldProps>;
};

export interface UpdateProps {
  update: React.Dispatch<React.SetStateAction<number>>;
}

export interface CreateUpdateProps extends CreateProps, UpdateProps {};

export interface EditUpdateProps extends EditProps, UpdateProps {};
