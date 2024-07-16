"use client";
import React, { useState } from "react";
import { Button, Chip } from "@nextui-org/react";
import { FaPlus, FaQrcode, FaEdit } from "react-icons/fa";
import { motion, AnimatePresence } from "framer-motion";
import AddOtpModal from "@/components/add-otp-modal";
import QRScanner from "@/components/qr-scanner";

const FAB = () => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [showQRScanner, setShowQRScanner] = useState(false);

  const toggleExpand = () => setIsExpanded(!isExpanded);

  return (
    <>
      <div className="fixed bottom-5 sm:bottom-10 right-5 sm:right-10 flex flex-col items-end space-y-4">
        <AnimatePresence>
          {isExpanded && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: 20 }}
              className="flex items-center space-x-2"
            >
              <Chip
                className="text-sm text-default-500"
                variant="solid"
                radius="sm"
              >
                Scan QR Code
              </Chip>
              <Button
                isIconOnly
                size="lg"
                variant="shadow"
                className="rounded-full w-14 h-14"
                onPress={() => setShowQRScanner(true)}
              >
                <FaQrcode />
              </Button>
            </motion.div>
          )}
          {isExpanded && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: 20 }}
              className="flex items-center space-x-2"
            >
              <Chip
                className="text-sm text-default-500"
                variant="solid"
                radius="sm"
              >
                Add Manually
              </Chip>
              <Button
                isIconOnly
                size="lg"
                className="rounded-full w-14 h-14"
                variant="shadow"
                onPress={() => setShowModal(true)}
              >
                <FaEdit />
              </Button>
            </motion.div>
          )}
        </AnimatePresence>
        <Button
          isIconOnly
          size="lg"
          className="rounded-full transition-transform transform hover:rotate-45 w-14 h-14"
          variant="shadow"
          color="success"
          onPress={toggleExpand}
        >
          <FaPlus className={isExpanded ? "rotate-45" : ""} />
        </Button>
      </div>

      <QRScanner
        isOpen={showQRScanner}
        onClose={() => setShowQRScanner(false)}
      />
      <AddOtpModal isOpen={showModal} onClose={() => setShowModal(false)} />
    </>
  );
};

export default FAB;
