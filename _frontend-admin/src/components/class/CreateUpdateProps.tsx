import { CreateProps } from 'react-admin';

interface UpdateProps {
  update: React.Dispatch<React.SetStateAction<number>>;
}

interface CreateUpdateProps extends CreateProps, UpdateProps {};

export default CreateUpdateProps;
