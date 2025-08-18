import { React, useState, useEffect, useRef } from 'react';
import './App.css';

import { UpdatePortList } from '../wailsjs/go/main/App'

import { theme, ConfigProvider, Flex, Layout, Table, Button, Switch, notification } from 'antd';

const { Sider, Content } = Layout;

function App() {
    const [portListData, setPortListData] = useState([]);
    const [autoScan, setAutoScan] = useState(false);
    const [backendError, setBackendError] = useState(null);

    const waitBackend = useRef(false);
    const autoScanIntervalId = useRef(null);

    useEffect(() => {
        if (autoScan && autoScanIntervalId.current === null)
            autoScanIntervalId.current = setInterval(updatePortList, 1000);

        if (!autoScan && autoScanIntervalId.current !== null) {
            clearInterval(autoScanIntervalId.current);
            autoScanIntervalId.current = null;
        }

        return () => {
            clearInterval(autoScanIntervalId.current);
            autoScanIntervalId.current = null;
        };
    }, [autoScan]);

    useEffect(() => {
        if (backendError) {
            notification.error({
                message: "Ошибка",
                description: backendError,
                placement: "bottomLeft"
            });
            setBackendError(null);
        }
    }, [backendError]);

    async function updatePortList() {
        try {
            console.log("updatePortList");
            console.log(`waitBackend: ${waitBackend.current}`);

            // Если бэк ещё не ответил на первый запрос,
            // то игонрируем новые вызовы с фронта (как ручные, так и по таймеру).
            // То есть, обеспечиваем синхронное взаимодействие между фронтом и бэком.
            if (waitBackend.current)
                return;

            waitBackend.current = true;
            let res = await UpdatePortList();
            waitBackend.current = false;

            let portList = res.map((portInfo, index) => ({ ...portInfo, key: index }));
            setPortListData(portList);
            console.table(portList);
        } catch (err) {
            console.log(`updatePortList: ${err}`);

            waitBackend.current = false;
            setBackendError(err);
        }
    }

    const portListCols = [
        { key: 0, title: 'Название', dataIndex: 'Name' },
        { key: 1, title: 'Usb', dataIndex: 'Usb' },
        { key: 2, title: 'VID', dataIndex: 'Vid' },
        { key: 3, title: 'PID', dataIndex: 'Pid' },
        { key: 4, title: 'Sent', dataIndex: 'SentData' },
        { key: 5, title: 'Received', dataIndex: 'ReceivedData' }
    ];

    const layoutStyle = {
        overflow: 'hidden',
        height: "100vh",
    };

    const siderStyle = {
        textAlign: 'left',
        lineHeight: '60px',
        color: '#fff',
        backgroundColor: '#ffffffff',
        padding: '16px',
        borderRight: "1px solid #f0f0f0ff"
    };

    const contentStyle = {
        textAlign: 'center',
        lineHeight: '120px',
        backgroundColor: '#ffffffff',
        display: "flex",
        flexDirection: "column"
    };

    return (
        <div id="App">
            <ConfigProvider
                theme={{
                    algorithm: theme.defaultAlgorithm,
                    token: {
                        borderRadius: 0,
                        fontFamily: "Roboto",
                    }
                }}
            >

                <Flex gap="middle" wrap>
                    <Layout style={layoutStyle}>
                        <Sider width="20%" style={siderStyle}>
                            <div
                                style={{
                                    display: "flex",
                                    flexDirection: "column",
                                    alignItems: "flex-start",
                                    justifyContent: "flex-start",
                                    gap: "0px",
                                    padding: "0px",
                                    height: "100%",
                                }}
                            >
                                <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                                    <Switch onChange={setAutoScan} />
                                    <span style={{ color: "black" }}>Автоскан</span>
                                </div>
                                <Button type="primary" onClick={updatePortList}>Обновить</Button>
                            </div>
                        </Sider>

                        <Content style={contentStyle}>
                            <Table
                                columns={portListCols}
                                dataSource={portListData}
                                pagination={false}
                                size="small"
                                style={{ flex: 1 }}
                            />
                        </Content>
                    </Layout>
                </Flex>
            </ConfigProvider>
        </div >
    )
}

export default App
