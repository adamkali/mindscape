import AdminRoute from './AdminRoute';
import CreateFolderComponent from './CreateFolderComponent';
import FolderComponent from './FolderComponent';
import { Header } from './Header';
import ProtectedRoute from './ProtectedRoute';
import * as atoms from './atoms';

const Components = {
	AdminRoute,
	CreateFolderComponent,
	FolderComponent,
	Header,
	ProtectedRoute,

	atoms: atoms,
}

export default Components;
