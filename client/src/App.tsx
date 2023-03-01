import { Button, Container, Flex, Input, ScrollArea } from '@mantine/core';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { v4 as uuid } from 'uuid';

import { WEBSOCKET_URL } from './utils/constants';

type socketMessage = {
    id: string;
    msg: string;
};

export const App = () => {
    const socket = useMemo(() => new WebSocket(WEBSOCKET_URL), [WEBSOCKET_URL]);
    const id = uuid();

    const [currMsg, setCurrMsg] = useState('');
    const [messages, setMessages] = useState<socketMessage[]>([]);

    const viewport = useRef<HTMLDivElement>(null);
    const scrollToBottom = useCallback(
        () =>
            viewport.current!.scrollTo({
                top: viewport.current!.scrollHeight,
                behavior: 'smooth',
            }),
        [viewport]
    );

    // const socketHello = useCallback(() => {
    //     const payload = {
    //         id,
    //         msg: `Hello from client ${id}!`,
    //     };
    //     socket.send(JSON.stringify(payload));
    // }, [socket]);

    const socketMessage = useCallback(async (ev: MessageEvent<any>) => {
        const payload = await JSON.parse(ev.data);
        setMessages((prev) => {
            return [...prev, { id: payload.id, msg: payload.msg }];
        });
        scrollToBottom();
    }, []);

    useEffect(() => {
        // socket.onopen = socketHello;
        socket.onmessage = socketMessage;
    }, []);

    return (
        <>
            <Container>
                <ScrollArea
                    w={'100%'}
                    style={{ height: 800 }}
                    viewportRef={viewport}
                >
                    {messages.map((msg, idx) => {
                        return <p key={idx}>{msg.msg}</p>;
                    })}
                </ScrollArea>

                <Flex
                    miw={'100%'}
                    justify='space-between'
                    align='flex-start'
                    direction='row'
                >
                    <Input
                        placeholder='Message'
                        w={'100%'}
                        mr={'md'}
                        value={currMsg}
                        onChange={(ev) => setCurrMsg(ev.target.value)}
                    />
                    <Button
                        ml={'md'}
                        onClick={() => {
                            const payload = {
                                id,
                                msg: currMsg,
                            };
                            socket.send(JSON.stringify(payload));
                            scrollToBottom();
                        }}
                    >
                        Send
                    </Button>
                </Flex>
            </Container>
        </>
    );
};
