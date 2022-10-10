import { Pagination } from 'react-admin';

const GlobalPagination = () => <Pagination rowsPerPageOptions={[ 25, 50, 100 ]} perPage={50} />;

export default GlobalPagination;
