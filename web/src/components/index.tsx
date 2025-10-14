import AdminRoute from './AdminRoute';
import CreateFolderComponent from './CreateFolderComponent';
import CreateBookmarkComponent from './CreateBookmarkComponent';
import BookmarkComponent from './BookmarkComponent';
import BookmarkCard from './BookmarkCard';
import FolderComponent from './FolderComponent';
import FolderCard from './FolderCard';
import { Header } from './Header';
import ProtectedRoute from './ProtectedRoute';
import * as atoms from './atoms';

const Components = {
	AdminRoute,
	CreateFolderComponent,
	FolderComponent,
	FolderCard,
	Header,
	ProtectedRoute,
	CreateBookmarkComponent,
	BookmarkComponent,
	BookmarkCard,

	atoms: atoms,
};

export default Components;
