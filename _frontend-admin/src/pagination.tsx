import { Pagination } from "react-admin";

export const GlobalPagination = () => <Pagination defaultValue={25} rowsPerPageOptions={[ 25, 50, 100 ]} />;
