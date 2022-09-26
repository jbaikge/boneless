import { Pagination } from "react-admin";

export const GlobalPagination = () => <Pagination perPage={50} rowsPerPageOptions={[ 25, 50, 100 ]} />;
