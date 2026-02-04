import AdminRoute from './AdminRoute';
import * as atoms from './atoms';
import BookmarkCard from './BookmarkCard';
import BookmarkComponent from './BookmarkComponent';
import CreateBookmarkComponent from './CreateBookmarkComponent';
import CreateFolderComponent from './CreateFolderComponent';
import FolderCard from './FolderCard';
import FolderComponent from './FolderComponent';
import { Header } from './Header';
import ProtectedRoute from './ProtectedRoute';
import WidgetContainer from './WidgetContainer';

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
	WidgetContainer,

	atoms: atoms,
};

export default Components;
