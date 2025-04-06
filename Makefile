
REACT_APP_PATH=es-project-react-app
REACT_APP_BUILD_DIR_PATH=${REACT_APP_PATH}/build

GO_WEB_DIR_PATH=web/

compile_react:

	cd ${REACT_APP_PATH} && npm run build;
	rm -rf ${GO_WEB_DIR_PATH}
	mv ${REACT_APP_BUILD_DIR_PATH} ${GO_WEB_DIR_PATH}

	
