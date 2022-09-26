import { Pagination } from "react-admin";

export const GlobalPagination = () => <Pagination rowsPerPageOptions={[ 25, 50, 100 ]} perPage={50} />;
