import AdminRoute from './AdminRoute';
import CreateFolderComponent from './CreateFolderComponent';
import CreateBookmarkComponent from './CreateBookmarkComponent';
import BookmarkComponent from './BookmarkComponent';
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
	CreateBookmarkComponent,
	BookmarkComponent,

	atoms: atoms,
}

export default Components;
